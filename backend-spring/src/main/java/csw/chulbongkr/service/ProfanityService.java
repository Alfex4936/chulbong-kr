package csw.chulbongkr.service;

import com.hankcs.algorithm.AhoCorasickDoubleArrayTrie;
import csw.chulbongkr.util.DAT;
import jakarta.annotation.PostConstruct;
import lombok.Getter;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;
import java.nio.charset.StandardCharsets;
import java.util.List;
import java.util.Objects;
import java.util.TreeMap;

@Slf4j
@Service
public class ProfanityService {
    private DAT<String> trie;
    @Getter
    private List<String> badWords;

    @PostConstruct
    public void init() throws IOException {
        trie = new DAT<>();
        TreeMap<String, String> badWordsMap = new TreeMap<>();

        // Load bad words from file
        try (BufferedReader reader = new BufferedReader(new InputStreamReader(
                Objects.requireNonNull(getClass().getResourceAsStream("/badwords.txt")), StandardCharsets.UTF_8))) {
            badWords = reader.lines()
                    .filter(word -> !word.isEmpty())
                    .distinct().toList();
            for (String word : badWords) {
                badWordsMap.put(word, word);
            }
        }

        trie.build(badWordsMap); // thread-safe after
    }

    public boolean containsProfanity(String text) {
        if (trie == null) {
            throw new IllegalStateException("ProfanityService not initialized");
        }
//        List<AhoCorasickDoubleArrayTrie.Hit<String>> hits = trie.parseText(text);
        DAT.Hit<String> hit = trie.findFirst(text);
        if (hit == null) {
            return false;
        }

        return hit.begin >= 0 && hit.end >= 0;
    }

}
