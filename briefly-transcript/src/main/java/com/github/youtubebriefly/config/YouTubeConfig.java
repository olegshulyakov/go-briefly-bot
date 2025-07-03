package com.github.youtubebriefly.config;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Configuration;
import org.springframework.util.StringUtils;

/**
 * Configuration class for YouTube-related proxy settings.
 * This class reads proxy configuration from application properties and constructs
 * a proxy string suitable for yt-dlp or similar tools.
 */
@Configuration
public class YouTubeConfig {
    @Value("${YOUTUBE_PROXY_HOST:}")
    private String youtubeProxyHost;

    @Value("${YOUTUBE_PROXY_PORT:}")
    private String youtubeProxyPort;

    @Value("${YOUTUBE_PROXY_USER:}")
    private String youtubeProxyUser;

    @Value("${YOUTUBE_PROXY_PASS:}")
    private String youtubeProxyPass;

    @Value("${YT_DLP_ADDITIONAL_OPTIONS:}")
    private String ytDlpAdditionalOptions;

    /**
     * Constructs the yt-dlp proxy string based on the configured properties.
     *
     * @return The proxy string (e.g., "socks5://user:pass@host:port") if configured,
     *         or an empty string if no proxy is configured, or null if configuration is incomplete.
     */
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
        return null;
    }

    /**
     * Gets the yt-dlp configured properties.
     *
     * @return The additional options if configured,
     *         or an empty string if no options is configured.
     */
    public String[] getYtDlpAdditionalOptions() {
        return StringUtils.hasText(this.ytDlpAdditionalOptions) ? this.ytDlpAdditionalOptions.split("\\s") : null;
    }
}
