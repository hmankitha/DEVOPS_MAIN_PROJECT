package com.example.ordermanagement.domain;

import jakarta.persistence.*;
import java.time.Instant;

@Entity
@Table(name = "shipments")
public class Shipment {
  @Id
  @GeneratedValue(strategy = GenerationType.IDENTITY)
  private Long id;

  @OneToOne(fetch = FetchType.LAZY)
  @JoinColumn(name = "order_id", nullable = false)
  private Order order;

  private String carrier; // e.g. UPS, DHL
  private String trackingNumber;

  @Enumerated(EnumType.STRING)
  private ShipmentStatus status = ShipmentStatus.PENDING;

  private Instant createdAt = Instant.now();
  private Instant shippedAt;
  private Instant deliveredAt;

  public Long getId() { return id; }
  public Order getOrder() { return order; }
  public void setOrder(Order order) { this.order = order; }
  public String getCarrier() { return carrier; }
  public void setCarrier(String carrier) { this.carrier = carrier; }
  public String getTrackingNumber() { return trackingNumber; }
  public void setTrackingNumber(String trackingNumber) { this.trackingNumber = trackingNumber; }
  public ShipmentStatus getStatus() { return status; }
  public void setStatus(ShipmentStatus status) { this.status = status; }
  public Instant getCreatedAt() { return createdAt; }
  public Instant getShippedAt() { return shippedAt; }
  public void setShippedAt(Instant shippedAt) { this.shippedAt = shippedAt; }
  public Instant getDeliveredAt() { return deliveredAt; }
  public void setDeliveredAt(Instant deliveredAt) { this.deliveredAt = deliveredAt; }
}
