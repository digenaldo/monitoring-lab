package com.monitoring.service;

import com.monitoring.model.Event;
import com.monitoring.repository.EventRepository;
import io.micrometer.core.instrument.Counter;
import io.micrometer.core.instrument.MeterRegistry;
import io.micrometer.core.instrument.Timer;
import org.springframework.data.mongodb.core.MongoTemplate;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;

@Service
public class MongoService {

    private final EventRepository eventRepository;
    private final MongoTemplate mongoTemplate;
    private final Counter operationsCounter;
    private final Timer latencyTimer;

    public MongoService(EventRepository eventRepository, 
                       MongoTemplate mongoTemplate,
                       MeterRegistry meterRegistry) {
        this.eventRepository = eventRepository;
        this.mongoTemplate = mongoTemplate;
        
        this.operationsCounter = Counter.builder("mongodb_operations_total")
                .description("Total number of MongoDB operations")
                .register(meterRegistry);
        
        this.latencyTimer = Timer.builder("mongodb_operation_duration_seconds")
                .description("MongoDB operation latency in seconds")
                .publishPercentiles(0.95, 0.99) // Publish P95 and P99 percentiles
                .register(meterRegistry);
    }

    public void performPeriodicOperations() {
        Timer.Sample sample = Timer.start();
        
        try {
            // Count operation
            mongoTemplate.getCollection("events").countDocuments();
            operationsCounter.increment();
            
            // Insert operation
            Event event = new Event();
            event.setTimestamp(LocalDateTime.now());
            event.setSource("spring-app");
            event.setMessage("Periodic event");
            eventRepository.save(event);
            operationsCounter.increment();
            
        } finally {
            sample.stop(latencyTimer);
        }
    }
}

