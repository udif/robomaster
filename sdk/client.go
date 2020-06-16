package sdk

import (
	"fmt"
	"github.com/brunoga/robomaster/sdk/modules/chassis"
	"net"

	"github.com/brunoga/robomaster/sdk/modules"
)

// Client enables controlling a RoboMaster robot through the plain-text SDK
// API (https://robomaster-dev.readthedocs.io/en/latest/).
type Client struct {
	finderModule  *modules.Finder
	controlModule *modules.Control
	pushModule    *modules.Push

	robotModule   *modules.Robot
	gimbalModule  *modules.Gimbal
	chassisModule *chassis.Chassis
	videoModule   *modules.Video
}

// NewClient returns a new client instance associated with the given ip. If ip
// is nil, the Client will try to detect a robot broadcasting its ip in the
// network.
func NewClient(ip net.IP) *Client {
	finderModule := modules.NewFinder()
	if ip != nil {
		finderModule.SetIP(ip)
	}

	controlModule := modules.NewControl(finderModule, false)
	pushModule := modules.NewPush(controlModule)
	robotModule := modules.NewRobot(controlModule)
	gimbalModule := modules.NewGimbal(controlModule, pushModule)
	chassisModule := chassis.New(controlModule, pushModule)
	videoModule := modules.NewVideo(controlModule)

	return &Client{
		finderModule,
		controlModule,
		pushModule,
		robotModule,
		gimbalModule,
		chassisModule,
		videoModule,
	}
}

// NewClientUSB creates a Client that tries to connect to the default USB
// connection ip.
func NewClientUSB() *Client {
	return NewClient(net.ParseIP("192.168.42.2"))
}

// NewClientWifiDirect creates a Client that tries to connect to the default
// WiFi Direct connection ip.
func NewClientWifiDirect() *Client {
	return NewClient(net.ParseIP("192.168.2.1"))
}

// Open opens the Client connection to the robot and enters SDK mode. Returns a
// nil error on success and a non-nil error on failure.
func (c *Client) Open() error {
	err := c.controlModule.Open()
	if err != nil {
		return fmt.Errorf("error starting control module: %w", err)
	}

	// Enter SDK mode.
	err = c.controlModule.SendDataExpectOk("command;")
	if err != nil {
		return fmt.Errorf("error entering to sdk mode: %w", err)
	}

	return nil
}

// Close exits SDk mode and closes the Client connection to the robot. Returns
// a nil error on success and a non-nil error on failure.
func (c *Client) Close() error {
	// Leave SDK mode.
	err := c.controlModule.SendDataExpectOk("quit;")
	if err != nil {
		return fmt.Errorf("error leaving sdk mode: %w", err)
	}

	err = c.controlModule.Close()
	if err != nil {
		return fmt.Errorf("error stopping control module")
	}

	return nil
}

// RobotModule returns a pointer to the associated Robot module. Used for
// doing generic robot-related operations.
func (c *Client) RobotModule() *modules.Robot {
	return c.robotModule
}

// GimbalModule returns a pointer to the associated Gimbal module. Used for
// doing gimbal-related operations.
func (c *Client) GimbalModule() *modules.Gimbal {
	return c.gimbalModule
}

// ChassisModule returns a pointer to the associated Chassis module. Used for
// doing chassis-related operations.
func (c *Client) ChassisModule() *chassis.Chassis {
	return c.chassisModule
}

// VideoModule returns a pointer to the associated Video module. Used for
// doing video-related operations.
func (c *Client) VideoModule() *modules.Video {
	return c.videoModule
}
