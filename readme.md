# BUILD GO-CODE FOR RASPBERRY PI
> GOARM=6 GOARCH=arm GOOS=linux go build -o=pi-uninacht .

# INSTALL FONTS ON RASPBERRY (MAKE SURE DIRECTORY ~/.fonts EXISTS ON PI)
> scp fonts/Open_Sans/* pi@192.168.2.22:~/.fonts
> scp fonts/Source_Code_Pro/* pi@192.168.2.22:~/.fonts

# INSTALL AND MANUALLY PUT FIREFOX (ICEWEASEL) IN FULLSCREEN MODE
> sudo apt-get update
> sudo apt-get install iceweasel

# COPY TO RASPBERRY PI (GIVEN PI HAS STATIC ADDRESS of 192.168.2.22)
# (MAKE SURE THAT DIRECTORIES /home/pi/uninacht AND /home/pi/uninacht/website EXIST)
> scp pi-uninacht pi@192.168.2.22:/home/pi/uninacht
> scp config.toml pi@192.168.2.22:/home/pi/uninacht
> scp run.sh pi@192.168.2.22:/home/pi/uninacht
> scp ./website/* pi@192.168.2.22:/home/pi/uninacht/website

# CHANGE PERMISSION ON RUN SCRIPT TO MAKE IT EXECUTABLE
> ssh pi@192.168.2.22
> sudo chmod +x /home/pi/uninacht/run.sh

# INSTALL AND MANUALLY PUT FIREFOX (ICEWEASEL) IN FULLSCREEN MODE
> sudo apt-get update
> sudo apt-get install iceweasel

# START UNINACHT-PROGRAM IMMEDIATELY AFTER STARTUP
> ssh pi@192.168.2.22
> sudo nano /etc/xdg/lxsession/LXDE-pi/autostart

enter the following line BEFORE the line with @xscreensaver:
@lxterminal --command "/home/pi/uninacht/run.sh"


# NOTE: TO RUN ICEWEASEL VIA SSH ON THE PI-MONITOR
> iceweasel --display=:0 http://localhost:8484 --fullscreen

# PHYSICAL CONNECTIONS
see pin assignment on http://pi.gadgetoid.com/pinout
use the physical numbers (in black) to match gobot gpio pins
