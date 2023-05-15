# Go Discord Token Grabber

<img alt="Code size" src="https://img.shields.io/github/languages/code-size/jonasrdl/gdtg?style=flat-square" > <img alt="Open Github issues" src="https://img.shields.io/github/issues/jonasrdl/gdtg?style=flat-square" >   

## Notice
**This project is made for educational purposes only.**

**Never give this token to third parties, the token can be used to perform any actions through your account.**
## Setup:
`git clone https://github.com/jonasrdl/gdtg`   

`go build gdtg.go`   

`./gdtg`

Available options to search:

`gdtg search all` Search in every location (Discord, Browsers, ...)   
Available options:
- Discord
- Discord Canary
- Google Chrome
- Brave (_stable_ and _nightly_ builds)

`gdtg search Discord` Search for exact string, for example "Discord"

`gdtg search /home/user/.config/discord` Search for custom path