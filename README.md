# Cron man

I'm writing this as my first project in my journey learning golang. I'm trying to write a CLI app to read list of cron job schedules from the program config file then create a daemon to run on the system to execute. Maybe later will add a terminal interface

## TODOS

    - [x] create config file (with viper)

    - [x] parse config file as map of Job

        - [] create default yaml with instructions file when first run

        - [] handle error for invalid format, empty config

    - [] support `cmdFile` for defining jobs

    - [] test execute with CLI arguments

    - [] add more unit test for executing job

    - [x] run in daemon? (reference: https://ieftimov.com/posts/four-steps-daemonize-your-golang-programs). Revert back, seems like I can run as daemon with launchd or supervisor

    - [] Logging solution

    - [] complete MVP

    - [] refactor into packages

    - [] Terminal User interface (bubbletea?)
