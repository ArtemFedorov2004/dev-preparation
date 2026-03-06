package com.devprep.controller;

import com.devprep.client.ApiClient;
import com.devprep.client.dto.Level;
import com.devprep.client.dto.TopicDto;
import com.devprep.client.dto.TopicWithQuestionsDto;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;

import java.util.List;
import java.util.Optional;

@Slf4j
@Controller
@RequestMapping("/topics")
@RequiredArgsConstructor
public class TopicController {

    private static final String PAGE_TITLE_TOPICS = "Темы — DevPrep";

    private static final String PAGE_TITLE_SUFFIX = " — DevPrep";

    private final ApiClient apiClient;

    @GetMapping
    public String listTopics(Model model) {
        List<TopicDto> topics = apiClient.getTopics();
        model.addAttribute("topics", topics);
        model.addAttribute("pageTitle", PAGE_TITLE_TOPICS);
        return "topics";
    }

    @GetMapping("/{slug}")
    public String topicDetail(
            @PathVariable String slug,
            @RequestParam(required = false, name = "level") String levelParam,
            Model model) {
        Level level = null;
        if (levelParam != null) {
            try {
                level = Level.fromValue(levelParam);
            } catch (IllegalArgumentException e) {
                log.warn("Invalid level parameter '{}'. Using default (null). Error: {}",
                        levelParam, e.getMessage());
            }
        }

        Optional<TopicWithQuestionsDto> topicOpt = apiClient.getTopicBySlug(slug, level);
        if (topicOpt.isEmpty()) {
            return "redirect:/topics";
        }
        TopicWithQuestionsDto topic = topicOpt.get();
        model.addAttribute("topic", topic);
        model.addAttribute("questions", topic.getQuestions() != null ? topic.getQuestions() : List.of());
        model.addAttribute("selectedLevel", level);
        model.addAttribute("levels", Level.values());
        model.addAttribute("pageTitle", topic.getName() + PAGE_TITLE_SUFFIX);
        return "topic-detail";
    }
}
