# dwarfpool-gomonitor
Dwarfpool Monero Simple Mine Monitor Tool


The purpose of this tool is to monitor your Dwarfpool Monero Mining Environment.  It works by taking in user defined paramaeters
and then queries Dwarfpool and extracts interesting information such as Hashrate, earned Monero coins, and estimated dollar amount.

Additionally, it has two primary modes being stand-alone and web-enabled.  Stand-alone simply returns the desired information on
the command line while web-enabled setups up a persistent web server.  Furthermore, in web-enabled mode the tool will automatically
notify the user by email when the hashate drops below a hard coded value.

Out of the box the tool will work on most systems (tested on windows 7 and Linux) under stand-alone mode.

The serv-noemail binary has been built for users who do not want email support but does want web-enabled mode support.  This binary also should work out of the box for most systems.

When used in web-enabled mode (for the serv binary) or if you compile yourself there is some configuration required.  The two primary requirements are setting up the correct environmental variables for the email service to work and depending on your gmail service you may need to tweak your account to allow the emails to bypass security settings.

On Linux:
- TOEMAIL="DestinationEmail"
- declare -x TOEMAIL
- FROMEMAIL="SenderEmail"
- declare -x FROMEMAIL
- PASSEMAIL="YourSenderEmailPassword"
- declare -x PASSEMAIL

On Windows:
- set TOEMAIL=DestinationEmail
- set FROMEMAIL=SenderEmail
- set PASSEMAIL=YourSenderEmailPassword


Currently, I also have additional code that monitors local CPU and GPU temperature and power on all client machines and performs some statistical calculations on the data.  However, it is a bit complex and rigid at the moment as it only works on Nvidia cards and was only tested for Linux machines.  If I have the cycles I will continue to work that or if anyone is really interested I can post that as well.



Happy Mining!
