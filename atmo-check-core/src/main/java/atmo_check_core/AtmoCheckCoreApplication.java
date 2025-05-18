package atmo_check_core;

import atmo_check_core.service.AirQualityService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.CommandLineRunner;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

@SpringBootApplication
public class AtmoCheckCoreApplication  implements CommandLineRunner {

	@Autowired
	private AirQualityService myService;

	public static void main(String[] args) {
		SpringApplication.run(AtmoCheckCoreApplication.class, args);
	}

	@Override
	public void run(String... args) throws Exception {
		myService.updateStations();
	}
}
