package com.example.ordermanagement.domain;

import jakarta.persistence.*;
import jakarta.validation.constraints.NotNull;
import java.math.BigDecimal;
import java.time.Instant;

@Entity
@Table(name = "payments")
public class Payment {
  @Id
  @GeneratedValue(strategy = GenerationType.IDENTITY)
  private Long id;

  @OneToOne(fetch = FetchType.LAZY)
  @JoinColumn(name = "order_id", nullable = false)
  private Order order;

  @NotNull
  private BigDecimal amount;

  @NotNull
  private String method; // e.g. CARD, UPI, PAYPAL

  @Enumerated(EnumType.STRING)
  private PaymentStatus status = PaymentStatus.INITIATED;

  private Instant createdAt = Instant.now();

  public Long getId() { return id; }
  public Order getOrder() { return order; }
  public void setOrder(Order order) { this.order = order; }
  public BigDecimal getAmount() { return amount; }
  public void setAmount(BigDecimal amount) { this.amount = amount; }
  public String getMethod() { return method; }
  public void setMethod(String method) { this.method = method; }
  public PaymentStatus getStatus() { return status; }
  public void setStatus(PaymentStatus status) { this.status = status; }
  public Instant getCreatedAt() { return createdAt; }
}
