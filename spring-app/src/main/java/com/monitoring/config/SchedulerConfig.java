package com.monitoring.config;

import com.monitoring.service.MongoService;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Component;

@Component
public class SchedulerConfig {

    private static final Logger logger = LoggerFactory.getLogger(SchedulerConfig.class);
    
    private final MongoService mongoService;

    public SchedulerConfig(MongoService mongoService) {
        this.mongoService = mongoService;
    }

    @Scheduled(fixedRate = 5000) // Every 5 seconds
    public void scheduledTask() {
        try {
            mongoService.performPeriodicOperations();
            logger.info("Periodic operations completed");
        } catch (Exception e) {
            logger.error("Error in scheduled task", e);
        }
    }
}

