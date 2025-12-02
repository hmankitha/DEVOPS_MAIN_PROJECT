package com.example.ordermanagement.api;

import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class HealthController {

    @GetMapping("/api/v1/orders/health")
    public String health() {
        return "{\"service\":\"order-management\",\"status\":\"healthy\"}";
    }
}