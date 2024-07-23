package csw.chulbongkr.service.search;

import org.apache.lucene.analysis.Analyzer;
import org.apache.lucene.analysis.TokenStream;
import org.apache.lucene.analysis.Tokenizer;
import org.apache.lucene.analysis.core.FlattenGraphFilter;
import org.apache.lucene.analysis.core.LowerCaseFilter;
import org.apache.lucene.analysis.miscellaneous.LengthFilter;
import org.apache.lucene.analysis.synonym.SynonymGraphFilter;
import org.apache.lucene.analysis.synonym.SynonymMap;
import org.apache.lucene.analysis.tokenattributes.CharTermAttribute;
import org.apache.lucene.analysis.tokenattributes.OffsetAttribute;
import org.apache.lucene.analysis.tokenattributes.PositionIncrementAttribute;
import org.apache.lucene.queryparser.classic.ParseException;
import org.apache.lucene.util.CharsRef;

import java.io.IOException;
import java.io.Reader;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class CustomKoreanAnalyzer extends Analyzer {

    private final SynonymMap synonymMap;

    public CustomKoreanAnalyzer() throws IOException, ParseException {
        synonymMap = buildSynonymMap();
    }

    @Override
    protected TokenStreamComponents createComponents(String fieldName) {
        Tokenizer source = new SimpleKoreanTokenizer();
        TokenStream filter = new LowerCaseFilter(source);
        filter = new LengthFilter(filter, 2, Integer.MAX_VALUE); // Ignore single characters
        filter = new SynonymGraphFilter(filter, synonymMap, true);
        filter = new FlattenGraphFilter(filter); // Flatten the graph for indexing
        return new TokenStreamComponents(source, filter);
    }

    private SynonymMap buildSynonymMap() throws IOException, ParseException {
        SynonymMap.Builder builder = new SynonymMap.Builder(true);

        // Define province synonyms
        Map<String, String> provinceMap = new HashMap<>();
        provinceMap.put("경기", "경기도");
        provinceMap.put("경기도", "경기도");
        provinceMap.put("ㄱㄱㄷ", "경기도");
        provinceMap.put("서울", "서울특별시");
        provinceMap.put("서울특별시", "서울특별시");
        provinceMap.put("ㅅㅇㅌㅂㅅ", "서울특별시");
        provinceMap.put("부산", "부산광역시");
        provinceMap.put("ㅄ", "부산광역시");
        provinceMap.put("ㅂㅅ", "부산광역시");
        provinceMap.put("ㅂㅅㄱㅇㅅ", "부산광역시");
        provinceMap.put("부산광역시", "부산광역시");
        provinceMap.put("대구", "대구광역시");
        provinceMap.put("대구광역시", "대구광역시");
        provinceMap.put("ㄷㄱㄱㅇㅅ", "대구광역시");
        provinceMap.put("인천", "인천광역시");
        provinceMap.put("인천광역시", "인천광역시");
        provinceMap.put("제주", "제주특별자치도");
        provinceMap.put("제주특별자치도", "제주특별자치도");
        provinceMap.put("ㅈㅈㄷ", "제주특별자치도");
        provinceMap.put("제주도", "제주특별자치도");
        provinceMap.put("대전", "대전광역시");
        provinceMap.put("대전광역시", "대전광역시");
        provinceMap.put("울산", "울산광역시");
        provinceMap.put("울산광역시", "울산광역시");
        provinceMap.put("광주", "광주광역시");
        provinceMap.put("광주광역시", "광주광역시");
        provinceMap.put("세종", "세종특별자치시");
        provinceMap.put("세종특별자치시", "세종특별자치시");
        provinceMap.put("ㅅㅈㅅ", "세종특별자치시");
        provinceMap.put("강원", "강원특별자치도");
        provinceMap.put("강원도", "강원특별자치도");
        provinceMap.put("강원특별자치도", "강원특별자치도");
        provinceMap.put("ㄱㅇㄷ", "강원특별자치도");
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

        for (Map.Entry<String, String> entry : provinceMap.entrySet()) {
            builder.add(new CharsRef(entry.getKey()), new CharsRef(entry.getValue()), true);
        }

        return builder.build();
    }


    public static class SimpleKoreanTokenizer extends Tokenizer {
        private final CharTermAttribute charTermAttribute = addAttribute(CharTermAttribute.class);
        private final OffsetAttribute offsetAttribute = addAttribute(OffsetAttribute.class);
        private final PositionIncrementAttribute positionIncrementAttribute = addAttribute(PositionIncrementAttribute.class);
        private final List<TokenInfo> tokens = new ArrayList<>();
        private int tokenIndex = 0;

        @Override
        public boolean incrementToken() throws IOException {
            if (tokenIndex < tokens.size()) {
                clearAttributes();
                TokenInfo tokenInfo = tokens.get(tokenIndex);
                charTermAttribute.append(tokenInfo.token);
                charTermAttribute.setLength(tokenInfo.token.length());
                offsetAttribute.setOffset(correctOffset(tokenInfo.startOffset), correctOffset(tokenInfo.endOffset));
                positionIncrementAttribute.setPositionIncrement(1);
                tokenIndex++;
                return true;
            }
            return false;
        }

        @Override
        public void reset() throws IOException {
            super.reset();
            tokenIndex = 0;
            tokens.clear();
            String text = inputToString(input);
            tokenize(text);
        }

        private String inputToString(Reader input) {
            StringBuilder builder = new StringBuilder();
            char[] buffer = new char[1024];
            int length;
            try {
                while ((length = input.read(buffer)) != -1) {
                    builder.append(buffer, 0, length);
                }
            } catch (IOException e) {
                throw new RuntimeException(e);
            }
            return builder.toString();
        }

        private void tokenize(String text) {
            int currentOffset = 0;
            for (String token : text.split("\\s+")) {
                int startOffset = text.indexOf(token, currentOffset);
                int endOffset = startOffset + token.length();
                tokens.add(new TokenInfo(token, startOffset, endOffset));
                currentOffset = endOffset;
            }
        }

        private static class TokenInfo {
            String token;
            int startOffset;
            int endOffset;

            TokenInfo(String token, int startOffset, int endOffset) {
                this.token = token;
                this.startOffset = startOffset;
                this.endOffset = endOffset;
            }
        }
    }
}
