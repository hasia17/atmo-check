package gios_data;

import org.springframework.boot.CommandLineRunner;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

@SpringBootApplication
public class AtmoCheckCoreApplication  implements CommandLineRunner {

	public static void main(String[] args) {
		SpringApplication.run(AtmoCheckCoreApplication.class, args);
	}

	@Override
	public void run(String... args) {
	}
}
