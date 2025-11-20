package com.example.ordermanagement.service;

import com.example.ordermanagement.api.dto.CreateOrderRequest;
import com.example.ordermanagement.api.dto.OrderItemDto;
import com.example.ordermanagement.domain.*;
import com.example.ordermanagement.repository.OrderRepository;
import io.micrometer.core.instrument.MeterRegistry;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.math.BigDecimal;

@Service
public class OrderService {
  private final OrderRepository orderRepository;
  private final MeterRegistry meterRegistry;

  public OrderService(OrderRepository orderRepository, MeterRegistry meterRegistry) {
    this.orderRepository = orderRepository;
    this.meterRegistry = meterRegistry;
  }

  @Transactional
  public Order createOrder(CreateOrderRequest req) {
    Order order = new Order();
    order.setCustomerEmail(req.getCustomerEmail());
    order.setStatus(OrderStatus.PENDING);

    BigDecimal total = BigDecimal.ZERO;
    for (OrderItemDto dto : req.getItems()) {
      OrderItem item = new OrderItem();
      item.setProductId(dto.getProductId());
      item.setQuantity(dto.getQuantity());
      item.setUnitPrice(dto.getUnitPrice());
      order.addItem(item);
      total = total.add(dto.getUnitPrice().multiply(BigDecimal.valueOf(dto.getQuantity())));
    }

    Payment payment = new Payment();
    payment.setAmount(total);
    payment.setMethod("CARD");
    order.setPayment(payment);

    meterRegistry.counter("orders.created").increment();
    return orderRepository.save(order);
  }

  @Transactional(readOnly = true)
  public Order getOrder(Long id) {
    return orderRepository.findById(id).orElse(null);
  }

  @Transactional
  public Order cancelOrder(Long id) {
    Order order = orderRepository.findById(id).orElse(null);
    if (order == null) return null;
    order.setStatus(OrderStatus.CANCELLED);
    meterRegistry.counter("orders.cancelled").increment();
    return orderRepository.save(order);
  }
}
