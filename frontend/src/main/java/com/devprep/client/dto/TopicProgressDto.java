package com.devprep.client.dto;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import lombok.Data;

@Data
@JsonIgnoreProperties(ignoreUnknown = true)
public class TopicProgressDto {

    private TopicDto topic;

    private int total;

    private int learned;

    private int needReview;

    private int dontKnow;
}