# DAQ-TFTP-Client

## Introduction
This application is a TFTP client in GUI form designed for RIT Launch Initiative members to be able to configure data acquisition boards.


## Usage
#### Step 1. Run the application with elevated privileges.

The client pings the given IP address before writing and receiving files which requires elevated privileges.

#### Step 2. Enter an IP address running a TFTP client.

#### Step 3. Open the instructions dropdown menu for more information.


## Installation
#### Step 1. Navigate to the [releases page](https://github.com/AarC10/DAQ-TFTP-Client/releases) for this repository.

#### Step 2. Download the executable for your own operating system and architecture. 

If your OS and architecture are not listed, let me know or navigate to the Building the Code section.
 
## Building the Code
#### Step 1. Install Go 1.17.6 or higher. The installation is linked [here](https://go.dev/dl/).

#### Step 2. Run ```go mod tidy```

#### Step 3. Run ```go build```

The executable designed for your system (unless specified otherwise in GOROOT) will be built in the the directory.
