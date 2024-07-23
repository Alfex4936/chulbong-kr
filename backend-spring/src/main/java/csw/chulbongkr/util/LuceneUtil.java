package csw.chulbongkr.util;

import lombok.extern.slf4j.Slf4j;

import java.util.HashMap;
import java.util.Map;
import java.util.Set;

import static java.util.Map.entry;

@Slf4j
public class LuceneUtil {
    private static final Set<Character> VALID_INITIAL_CONSONANTS = Set.of(
            'ㄱ', 'ㄲ', 'ㄴ', 'ㄷ', 'ㄸ', 'ㄹ', 'ㅁ', 'ㅂ', 'ㅃ', 'ㅅ', 'ㅆ', 'ㅇ', 'ㅈ', 'ㅉ', 'ㅊ', 'ㅋ', 'ㅌ', 'ㅍ', 'ㅎ'
    );

    private static final Map<Character, char[]> DOUBLE_CONSONANTS = Map.ofEntries(
            entry('ㄳ', new char[] {'ㄱ', 'ㅅ'}),
            entry('ㄵ', new char[] {'ㄴ', 'ㅈ'}),
            entry('ㄶ', new char[] {'ㄴ', 'ㅎ'}),
            entry('ㄺ', new char[] {'ㄹ', 'ㄱ'}),
            entry('ㄻ', new char[] {'ㄹ', 'ㅁ'}),
            entry('ㄼ', new char[] {'ㄹ', 'ㅂ'}),
            entry('ㄽ', new char[] {'ㄹ', 'ㅅ'}),
            entry('ㄾ', new char[] {'ㄹ', 'ㅌ'}),
            entry('ㄿ', new char[] {'ㄹ', 'ㅍ'}),
            entry('ㅀ', new char[] {'ㄹ', 'ㅎ'}),
            entry('ㅄ', new char[] {'ㅂ', 'ㅅ'})
    );

    private static final Map<Character, Character> INITIAL_CONSONANT_MAP = new HashMap<>();

    static {
        INITIAL_CONSONANT_MAP.put('\u1100', 'ㄱ');
        INITIAL_CONSONANT_MAP.put('\u1101', 'ㄲ');
        INITIAL_CONSONANT_MAP.put('\u1102', 'ㄴ');
        INITIAL_CONSONANT_MAP.put('\u1103', 'ㄷ');
        INITIAL_CONSONANT_MAP.put('\u1104', 'ㄸ');
        INITIAL_CONSONANT_MAP.put('\u1105', 'ㄹ');
        INITIAL_CONSONANT_MAP.put('\u1106', 'ㅁ');
        INITIAL_CONSONANT_MAP.put('\u1107', 'ㅂ');
        INITIAL_CONSONANT_MAP.put('\u1108', 'ㅃ');
        INITIAL_CONSONANT_MAP.put('\u1109', 'ㅅ');
        INITIAL_CONSONANT_MAP.put('\u110A', 'ㅆ');
        INITIAL_CONSONANT_MAP.put('\u110B', 'ㅇ');
        INITIAL_CONSONANT_MAP.put('\u110C', 'ㅈ');
        INITIAL_CONSONANT_MAP.put('\u110D', 'ㅉ');
        INITIAL_CONSONANT_MAP.put('\u110E', 'ㅊ');
        INITIAL_CONSONANT_MAP.put('\u110F', 'ㅋ');
        INITIAL_CONSONANT_MAP.put('\u1110', 'ㅌ');
        INITIAL_CONSONANT_MAP.put('\u1111', 'ㅍ');
        INITIAL_CONSONANT_MAP.put('\u1112', 'ㅎ');
    }

    private static final Map<String, String> provinceMap = new HashMap<>();

    static {
        provinceMap.put("경기", "경기도");
        provinceMap.put("경기도", "경기도");
        provinceMap.put("서울", "서울특별시");
        provinceMap.put("서울특별시", "서울특별시");
        provinceMap.put("부산", "부산광역시");
        provinceMap.put("부산광역시", "부산광역시");
        provinceMap.put("대구", "대구광역시");
        provinceMap.put("대구광역시", "대구광역시");
        provinceMap.put("인천", "인천광역시");
        provinceMap.put("인천광역시", "인천광역시");
        provinceMap.put("제주", "제주특별자치도");
        provinceMap.put("제주특별자치도", "제주특별자치도");
        provinceMap.put("제주도", "제주특별자치도");
        provinceMap.put("대전", "대전광역시");
        provinceMap.put("대전광역시", "대전광역시");
        provinceMap.put("울산", "울산광역시");
        provinceMap.put("울산광역시", "울산광역시");
        provinceMap.put("광주", "광주광역시");
        provinceMap.put("광주광역시", "광주광역시");
        provinceMap.put("세종", "세종특별자치시");
        provinceMap.put("세종특별자치시", "세종특별자치시");
        provinceMap.put("강원", "강원특별자치도");
        provinceMap.put("강원도", "강원특별자치도");
        provinceMap.put("강원특별자치도", "강원특별자치도");
        provinceMap.put("경남", "경상남도");
        provinceMap.put("경상남도", "경상남도");
        provinceMap.put("경북", "경상북도");
        provinceMap.put("경상북도", "경상북도");
        provinceMap.put("전북", "전라북도");
        provinceMap.put("전북특별자치도", "전라북도");
        provinceMap.put("충남", "충청남도");
        provinceMap.put("충청남도", "충청남도");
        provinceMap.put("충북", "충청북도");
        provinceMap.put("충청북도", "충청북도");
        provinceMap.put("전남", "전라남도");
        provinceMap.put("전라남도", "전라남도");
    }

    public static String standardizeProvince(String province) {
        return provinceMap.getOrDefault(province, province);
    }

    public static String extractInitialConsonants(String text) {
        StringBuilder initialConsonants = new StringBuilder();
        for (char ch : text.toCharArray()) {
            if (ch >= 0xAC00 && ch <= 0xD7A3) {
                int unicode = ch - 0xAC00;
                int choIndex = unicode / (21 * 28);
                char initialConsonant = (char) (0x1100 + choIndex);
                initialConsonants.append(INITIAL_CONSONANT_MAP.getOrDefault(initialConsonant, ch));
            } else if (INITIAL_CONSONANT_MAP.containsValue(ch)) {
                initialConsonants.append(ch);
            }
        }
        return initialConsonants.toString();
    }

    public static String segmentConsonants(String input) {
        var result = new StringBuilder();

        for (char ch : input.toCharArray()) {
            if (VALID_INITIAL_CONSONANTS.contains(ch)) {
                result.append(ch);
            } else if (DOUBLE_CONSONANTS.containsKey(ch)) {
                result.append(DOUBLE_CONSONANTS.get(ch));
            } else {
                result.append(ch);
            }
        }

        return result.toString();
    }
}
