package com.github.youtubebriefly.config;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.util.StringUtils;

@Configuration
public class YoutubeProxyConfig {
    @Value("${YOUTUBE_PROXY_HOST:}")
    private String youtubeProxyHost;

    @Value("${YOUTUBE_PROXY_PORT:}")
    private String youtubeProxyPort;

    @Value("${YOUTUBE_PROXY_USER:}")
    private String youtubeProxyUser;

    @Value("${YOUTUBE_PROXY_PASS:}")
    private String youtubeProxyPass;

    @Bean
    public String getYtDlpProxy() {
        if (StringUtils.hasText(this.youtubeProxyHost)
                && StringUtils.hasText(this.youtubeProxyPort)
                && StringUtils.hasText(this.youtubeProxyUser)
                && StringUtils.hasText(this.youtubeProxyPass)
        ) {
            return String.format("socks5://%s:%s@%s:%s", this.youtubeProxyUser, this.youtubeProxyPass, this.youtubeProxyHost, this.youtubeProxyPort);
        } else if (StringUtils.hasText(this.youtubeProxyHost)
                && StringUtils.hasText(this.youtubeProxyPort)
        ) {
            return String.format("socks5://%s:%s", this.youtubeProxyHost, this.youtubeProxyPort);
        }
        return "";
    }
}
