package com.devprep.controller;

import com.devprep.client.ApiClient;
import com.devprep.client.dto.BookmarkDto;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@Controller
@RequestMapping("/me")
@RequiredArgsConstructor
public class MyBookmarksController {

    private static final String PAGE_TITLE_BOOKMARKS = "Закладки — DevPrep";

    private final ApiClient apiClient;

    @GetMapping("/bookmarks")
    public String bookmarks(Model model) {
        setupBookmarksPageModel(model);
        return "me/bookmarks";
    }

    @PostMapping("/bookmarks")
    public String toggleBookmark(@RequestParam String slug, Model model) {
        apiClient.toggleBookmark(slug);
        setupBookmarksPageModel(model);
        return "redirect:/me/bookmarks";
    }

    @PostMapping("/questions/{slug}/bookmark")
    public String toggleBookmark(@PathVariable String slug) {
        apiClient.toggleBookmark(slug);
        return "redirect:/questions/%s".formatted(slug);
    }

    private void setupBookmarksPageModel(Model model) {
        List<BookmarkDto> bookmarks = apiClient.getMyBookmarks();
        model.addAttribute("bookmarks", bookmarks);
        model.addAttribute("pageTitle", PAGE_TITLE_BOOKMARKS);
    }
}
