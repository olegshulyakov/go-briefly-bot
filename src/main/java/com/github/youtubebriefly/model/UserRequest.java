package com.github.youtubebriefly.model;

import jakarta.persistence.Entity;
import jakarta.persistence.GeneratedValue;
import jakarta.persistence.Id;
import jakarta.persistence.Table;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.time.LocalDateTime;

@Entity
@Table(name = "UserRequest")
@Data
@AllArgsConstructor
@NoArgsConstructor
public class UserRequest {
    @Id
    @GeneratedValue
    private Long id;
    private Long userId; // Telegram user ID
    private String videoId;
    private LocalDateTime requestedAt;
}
