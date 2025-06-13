package com.github.youtubebriefly.dao;

import com.github.youtubebriefly.model.UserRequest;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface UserRequestRepository extends JpaRepository<UserRequest, Long> {
}
