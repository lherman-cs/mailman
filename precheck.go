package main

func init() {
	// TODO! Exit if either udp or tcp is already in use
	// _, err := net.DialTimeout("udp", fmt.Sprintf(":%d", Port), time.Second*2)
	// if err == nil {
	// 	log.Fatalf("udp port at %d is already in use", Port)
	// }

	// _, err = net.DialTimeout("tcp", fmt.Sprintf(":%d", Port), time.Second*2)
	// if err == nil {
	// 	log.Fatalf("tcp port at %d is already in use", Port)
	// }
}
