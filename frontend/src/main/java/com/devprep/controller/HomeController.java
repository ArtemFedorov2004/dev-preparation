package com.devprep.controller;

import com.devprep.client.ApiClient;
import com.devprep.client.dto.TopicDto;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.web.bind.annotation.GetMapping;

import java.util.List;

@Controller
@RequiredArgsConstructor
public class HomeController {

    private static final String PAGE_TITLE_HOME = "DevPrep — Платформа для подготовки к собеседованиям";

    private final ApiClient apiClient;

    @GetMapping("/")
    public String home(Model model) {
        List<TopicDto> topics = apiClient.getTopics();
        model.addAttribute("topics", topics);
        model.addAttribute("pageTitle", PAGE_TITLE_HOME);
        return "index";
    }
}
