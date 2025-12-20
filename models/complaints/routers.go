package complaints

import (
	"net/http"
	"strconv"

	"plantheon-backend/models/users"

	"github.com/gin-gonic/gin"
)

// CreateComplaintHandler handles creating a new complaint
func CreateComplaintHandler(c *gin.Context) {
	// Get user from context
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	user, ok := userInterface.(*users.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user format"})
		return
	}

	var req CreateComplaintRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if err := ValidateCreateComplaintRequest(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check for duplicate complaint
	exists, err := CheckDuplicateComplaint(user.ID, req.TargetID, ComplaintType(req.TargetType))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check duplicate complaint"})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "You have already submitted a complaint for this content"})
		return
	}

	complaint := &Complaint{
		UserID:     user.ID,
		TargetID:   req.TargetID,
		TargetType: ComplaintType(req.TargetType),
		Category:   ComplaintCategory(req.Category),
		Content:    req.Content,
		Status:     ComplaintStatusPending,
	}

	if err := CreateComplaint(complaint); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create complaint"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Complaint submitted successfully",
		"data":    complaint.ToComplaintResponse(),
	})
}

// CreateScanComplaintHandler handles creating a new complaint about scan results
func CreateScanComplaintHandler(c *gin.Context) {
	// Get user from context
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	user, ok := userInterface.(*users.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user format"})
		return
	}

	var req CreateScanComplaintRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if err := ValidateCreateScanComplaintRequest(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// No duplicate check for scan complaints - users can submit multiple complaints for the same disease

	complaint := &Complaint{
		UserID:                 user.ID,
		TargetID:               req.PredictedDiseaseID, // Use predicted disease as target
		TargetType:             ComplaintTypeScan,
		Category:               ComplaintCategory(req.Category),
		Content:                req.Content,
		ImageURL:               req.ImageURL,
		Status:                 ComplaintStatusPending,
		PredictedDiseaseID:     &req.PredictedDiseaseID,
		UserSuggestedDiseaseID: req.UserSuggestedDiseaseID,
		ConfidenceScore:        &req.ConfidenceScore,
		IsVerified:             false,
	}

	if err := CreateComplaint(complaint); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create complaint"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Scan complaint submitted successfully",
		"data":    complaint.ToComplaintResponse(),
	})
}

// GetMyComplaintsHandler gets all complaints made by the current user
func GetMyComplaintsHandler(c *gin.Context) {
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	user, ok := userInterface.(*users.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user format"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit

	complaints, total, err := GetComplaintsByUserID(user.ID, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get complaints"})
		return
	}

	responses := make([]ComplaintResponse, len(complaints))
	for i, complaint := range complaints {
		responses[i] = complaint.ToComplaintResponse()
	}

	c.JSON(http.StatusOK, gin.H{
		"data": ComplaintListResponse{
			Complaints: responses,
			Total:      total,
			Page:       page,
			Limit:      limit,
		},
	})
}

// GetComplaintByIDHandler gets a complaint by ID
func GetComplaintByIDHandler(c *gin.Context) {
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	user, ok := userInterface.(*users.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user format"})
		return
	}

	id := c.Param("id")
	if err := ValidateUUID(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	complaint, err := GetComplaintByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Complaint not found"})
		return
	}

	// Only allow user to see their own complaints (unless admin)
	if complaint.UserID != user.ID && user.Role != "ADMIN" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only view your own complaints"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": complaint.ToComplaintResponse()})
}

// GetComplaintsAboutMyContentHandler gets all complaints about the current user's posts/comments
func GetComplaintsAboutMyContentHandler(c *gin.Context) {
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	user, ok := userInterface.(*users.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user format"})
		return
	}

	status := c.Query("status")
	targetType := c.Query("target_type")

	complaints, err := GetComplaintsAboutMyContent(user.ID, status, targetType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get complaints"})
		return
	}

	responses := make([]ComplaintWithTargetResponse, len(complaints))
	for i, complaint := range complaints {
		responses[i] = complaint.ToComplaintWithTargetResponse()
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  responses,
		"total": len(responses),
	})
}

// DeleteComplaintHandler deletes a complaint (user can only delete their own pending complaints)
func DeleteComplaintHandler(c *gin.Context) {
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	user, ok := userInterface.(*users.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user format"})
		return
	}

	id := c.Param("id")
	if err := ValidateUUID(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	complaint, err := GetComplaintByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Complaint not found"})
		return
	}

	// Only allow user to delete their own pending complaints
	if complaint.UserID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own complaints"})
		return
	}

	if complaint.Status != ComplaintStatusPending {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You can only delete pending complaints"})
		return
	}

	if err := DeleteComplaint(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete complaint"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Complaint deleted successfully"})
}

// ============ ADMIN HANDLERS ============

// GetAllComplaintsHandler gets all complaints with pagination and filters (admin only)
func GetAllComplaintsHandler(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")
	targetType := c.Query("target_type")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit

	complaints, total, err := GetAllComplaints(offset, limit, status, targetType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get complaints"})
		return
	}

	responses := make([]ComplaintResponse, len(complaints))
	for i, complaint := range complaints {
		responses[i] = complaint.ToComplaintResponse()
	}

	c.JSON(http.StatusOK, gin.H{
		"data": ComplaintListResponse{
			Complaints: responses,
			Total:      total,
			Page:       page,
			Limit:      limit,
		},
	})
}

// GetComplaintsCountHandler gets the count of complaints with optional filters (admin only)
func GetComplaintsCountHandler(c *gin.Context) {
	status := c.Query("status")
	targetType := c.Query("target_type")
	isVerifiedStr := c.Query("is_verified")
	
	var isVerified *bool
	if isVerifiedStr != "" {
		val := isVerifiedStr == "true"
		isVerified = &val
	}

	count, err := GetComplaintsCount(status, targetType, isVerified)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get complaints count"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"count":       count,
			"status":      status,
			"target_type": targetType,
			"is_verified": isVerifiedStr,
		},
	})
}

// UpdateComplaintStatusHandler updates the status of a complaint (admin only)
func UpdateComplaintStatusHandler(c *gin.Context) {
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	user, ok := userInterface.(*users.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user format"})
		return
	}

	id := c.Param("id")
	if err := ValidateUUID(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req UpdateComplaintStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if err := ValidateUpdateComplaintStatusRequest(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if complaint exists
	_, err := GetComplaintByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Complaint not found"})
		return
	}

	if err := UpdateComplaintStatus(id, ComplaintStatus(req.Status), req.AdminNotes, user.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update complaint status"})
		return
	}

	// Get updated complaint
	updatedComplaint, _ := GetComplaintByID(id)

	c.JSON(http.StatusOK, gin.H{
		"message": "Complaint status updated successfully",
		"data":    updatedComplaint.ToComplaintResponse(),
	})
}

// AdminDeleteComplaintHandler deletes any complaint (admin only)
func AdminDeleteComplaintHandler(c *gin.Context) {
	id := c.Param("id")
	if err := ValidateUUID(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := GetComplaintByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Complaint not found"})
		return
	}

	if err := DeleteComplaint(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete complaint"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Complaint deleted successfully"})
}

// GetComplaintsByTargetHandler gets all complaints for a specific post, comment, or scan (admin only)
func GetComplaintsByTargetHandler(c *gin.Context) {
	targetID := c.Param("target_id")
	targetType := c.Query("target_type")

	if err := ValidateUUID(targetID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "target_id must be a valid UUID"})
		return
	}

	if targetType != string(ComplaintTypePost) && targetType != string(ComplaintTypeComment) && targetType != string(ComplaintTypeScan) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "target_type must be POST, COMMENT, or SCAN"})
		return
	}

	complaints, err := GetComplaintsByTargetID(targetID, ComplaintType(targetType))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get complaints"})
		return
	}

	responses := make([]ComplaintResponse, len(complaints))
	for i, complaint := range complaints {
		responses[i] = complaint.ToComplaintResponse()
	}

	c.JSON(http.StatusOK, gin.H{"data": responses})
}

// ============ SCAN COMPLAINT VERIFICATION HANDLERS ============

// VerifyComplaintHandler verifies a scan complaint and sets ground truth (admin only)
func VerifyComplaintHandler(c *gin.Context) {
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	user, ok := userInterface.(*users.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user format"})
		return
	}

	id := c.Param("id")
	if err := ValidateUUID(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req VerifyComplaintRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if err := ValidateVerifyComplaintRequest(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if complaint exists and is a scan complaint
	complaint, err := GetComplaintByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Complaint not found"})
		return
	}

	if complaint.TargetType != ComplaintTypeScan {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only scan complaints can be verified"})
		return
	}

	// Verify the complaint
	if err := VerifyComplaint(id, req.VerifiedDiseaseID, req.IsVerified, req.AdminNotes, user.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify complaint"})
		return
	}

	// Get updated complaint
	updatedComplaint, _ := GetComplaintByID(id)

	c.JSON(http.StatusOK, gin.H{
		"message": "Complaint verified successfully",
		"data":    updatedComplaint.ToComplaintResponse(),
	})
}

// GetUnverifiedScanComplaintsHandler gets unverified scan complaints (admin only)
func GetUnverifiedScanComplaintsHandler(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit

	complaints, total, err := GetUnverifiedScanComplaints(offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get unverified complaints"})
		return
	}

	responses := make([]ComplaintResponse, len(complaints))
	for i, complaint := range complaints {
		responses[i] = complaint.ToComplaintResponse()
	}

	c.JSON(http.StatusOK, gin.H{
		"data": ComplaintListResponse{
			Complaints: responses,
			Total:      total,
			Page:       page,
			Limit:      limit,
		},
	})
}

// ============ ML EXPORT HANDLER ============

// ExportTrainingDataHandler exports verified scan complaints for ML training (admin only)
func ExportTrainingDataHandler(c *gin.Context) {
	data, err := ExportTrainingData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to export training data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Training data exported successfully",
		"count":   len(data),
		"data":    data,
	})
}

// GetProblematicDiseasesHandler returns diseases with most complaints
func GetProblematicDiseasesHandler(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	data, err := GetMostProblematicDiseases(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get problematic diseases"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Problematic diseases retrieved successfully",
		"count":   len(data),
		"data":    data,
	})
}

// GetComplaintTrendsHandler returns daily complaint trends
func GetComplaintTrendsHandler(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))

	data, err := GetComplaintTrends(days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get complaint trends"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Complaint trends retrieved successfully",
		"days":    days,
		"count":   len(data),
		"data":    data,
	})
}

// GetOverallStatsHandler returns overall system statistics
func GetOverallStatsHandler(c *gin.Context) {
	data, err := GetOverallStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get overall stats"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Overall statistics retrieved successfully",
		"data":    data,
	})
}

// GetTopContributorsHandler returns users who contributed most
func GetTopContributorsHandler(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	data, err := GetTopContributors(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get top contributors"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Top contributors retrieved successfully",
		"count":   len(data),
		"data":    data,
	})
}
