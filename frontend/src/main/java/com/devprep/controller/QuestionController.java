package com.devprep.controller;

import com.devprep.client.ApiClient;
import com.devprep.client.dto.*;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.security.authentication.AnonymousAuthenticationToken;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.context.SecurityContextHolder;
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
@RequestMapping("/questions")
@RequiredArgsConstructor
public class QuestionController {

    private static final String PAGE_TITLE_QUESTIONS = "Вопросы — DevPrep";

    private static final String PAGE_TITLE_SUFFIX = " — DevPrep";

    private final ApiClient apiClient;

    @GetMapping
    public String listQuestions(
            @RequestParam(required = false) String topic,
            @RequestParam(required = false) String tag,
            @RequestParam(required = false, name = "level") String levelParam,
            @RequestParam(defaultValue = "1") int page,
            @RequestParam(defaultValue = "20") int limit,
            Model model) {
        Level level = null;
        if (levelParam != null && !levelParam.trim().isEmpty()) {
            try {
                level = Level.fromValue(levelParam);
            } catch (IllegalArgumentException e) {
                log.warn("Invalid level parameter '{}'. Using default (null).", levelParam);
            }
        }

        PaginatedQuestionsDto result = apiClient.getQuestions(topic, tag, level, page, limit);
        List<TopicDto> topics = apiClient.getTopics();
        List<TagDto> tags = apiClient.getTags();

        model.addAttribute("questions", result.getData());
        model.addAttribute("pagination", result.getPagination());
        model.addAttribute("topics", topics);
        model.addAttribute("tags", tags);
        model.addAttribute("levels", Level.values());
        model.addAttribute("selectedTopic", topic);
        model.addAttribute("selectedTag", tag);
        model.addAttribute("selectedLevel", level);
        model.addAttribute("pageTitle", PAGE_TITLE_QUESTIONS);
        return "questions";
    }

    @GetMapping("/{slug}")
    public String questionDetail(@PathVariable String slug, Model model) {
        Optional<QuestionDetailDto> questionOpt = apiClient.getQuestionBySlug(slug);
        if (questionOpt.isEmpty()) {
            return "redirect:/questions";
        }
        QuestionDetailDto question = questionOpt.get();
        model.addAttribute("question", question);
        model.addAttribute("pageTitle", question.getTitle() + PAGE_TITLE_SUFFIX);

        if (isAuthenticated()) {
            Optional<ProgressStatus> progressStatus = apiClient.getQuestionProgress(question.getSlug());
            boolean bookmarked = apiClient.isBookmarked(question.getSlug());
            model.addAttribute("progressStatus", progressStatus.orElse(null));
            model.addAttribute("bookmarked", bookmarked);
            model.addAttribute("progressStatuses", ProgressStatus.values());

            apiClient.recordView(question.getSlug());
        }

        return "question-detail";
    }

    private boolean isAuthenticated() {
        Authentication auth = SecurityContextHolder.getContext().getAuthentication();
        return auth != null
                && !(auth instanceof AnonymousAuthenticationToken)
                && auth.isAuthenticated();
    }
}