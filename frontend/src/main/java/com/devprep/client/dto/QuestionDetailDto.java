package com.devprep.client.dto;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import lombok.Data;

import java.util.List;

@Data
@JsonIgnoreProperties(ignoreUnknown = true)
public class QuestionDetailDto {

    private int id;

    private String slug;

    private String title;

    private String answer;

    private Level level;

    private TopicDto topic;

    private List<TagDto> tags;
}