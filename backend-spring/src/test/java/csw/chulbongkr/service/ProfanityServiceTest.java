package csw.chulbongkr.service;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.mockito.InjectMocks;
import org.mockito.MockitoAnnotations;

import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertTrue;

class ProfanityServiceTest {

    @InjectMocks
    private ProfanityService profanityService;

    @BeforeEach
    void setUp() throws Exception {
        MockitoAnnotations.openMocks(this);

        // Manually initialize the ProfanityService
        profanityService.init();
    }

    @Test
    void testContainsProfanity() {
        String textWithProfanity = "가나다라마사바아사아자차아카파타 시발이이임.";
        String textWithoutProfanity = "좋네요 아하하하하 가믈ㄴㄹㅇ12030ㅁㄴㅇ츸ㅌ9.";

        assertTrue(profanityService.containsProfanity(textWithProfanity));
        assertFalse(profanityService.containsProfanity(textWithoutProfanity));
    }

    @Test
    void testInitLoadsBadWords() {
        // This test assumes the bad words list is correctly loaded during initialization
        String text = "시발";
        assertTrue(profanityService.containsProfanity(text));

        // A text without any bad words
        String cleanText = "안녕";
        assertFalse(profanityService.containsProfanity(cleanText));
    }
}