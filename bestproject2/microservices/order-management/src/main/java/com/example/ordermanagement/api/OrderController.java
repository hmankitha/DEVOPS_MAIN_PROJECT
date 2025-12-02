package com.example.ordermanagement.api;

import com.example.ordermanagement.api.dto.CreateOrderRequest;
import com.example.ordermanagement.domain.Order;
import com.example.ordermanagement.service.OrderService;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.validation.Valid;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.net.URI;

@RestController
@RequestMapping("/api/v1/orders")
@Tag(name = "Orders", description = "Order management endpoints")
public class OrderController {

    private final OrderService orderService;

    public OrderController(OrderService orderService) {
        this.orderService = orderService;
    }

    @Operation(summary = "Create a new order")
    @PostMapping
    public ResponseEntity<Order> create(@Valid @RequestBody CreateOrderRequest req) {
        Order saved = orderService.createOrder(req);
        return ResponseEntity.created(URI.create("/api/v1/orders/" + saved.getId())).body(saved);
    }

    @Operation(summary = "Get order by id")
    // Only match numeric IDs — prevents conflict with /health
    @GetMapping("/{id:\\d+}")
    public ResponseEntity<Order> getOrder(@PathVariable("id") Long id) {
        Order o = orderService.getOrder(id);
        return (o == null) ? ResponseEntity.notFound().build() : ResponseEntity.ok(o);
    }

    @Operation(summary = "Cancel order")
    @PostMapping("/{id:\\d+}/cancel")
    public ResponseEntity<Order> cancel(@PathVariable Long id) {
        Order o = orderService.cancelOrder(id);
        return (o == null) ? ResponseEntity.notFound().build() : ResponseEntity.ok(o);
    }
}