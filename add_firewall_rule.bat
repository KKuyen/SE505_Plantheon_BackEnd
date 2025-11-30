@echo off
echo Adding firewall rule for Golang server on port 3000...
netsh advfirewall firewall add rule name="Golang Server" dir=in action=allow protocol=TCP localport=3000
echo Done! Firewall rule added.
pause
