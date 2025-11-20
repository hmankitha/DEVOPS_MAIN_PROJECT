package com.example.ordermanagement;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.boot.context.event.ApplicationReadyEvent;
import org.springframework.context.event.EventListener;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

@SpringBootApplication
public class OrderManagementApplication {
  private static final Logger log = LoggerFactory.getLogger(OrderManagementApplication.class);
  public static void main(String[] args) {
    SpringApplication.run(OrderManagementApplication.class, args);
  }

  @EventListener(ApplicationReadyEvent.class)
  public void ready() {
    log.info("Order Management Service started");
  }
}
