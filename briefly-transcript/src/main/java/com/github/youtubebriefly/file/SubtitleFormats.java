package com.github.youtubebriefly.file;

public enum SubtitleFormats {
    ASS,
    LRC,
    SRT,
    VTT;

    @Override
    public String toString() {
        return this.name().toLowerCase();
    }
}
