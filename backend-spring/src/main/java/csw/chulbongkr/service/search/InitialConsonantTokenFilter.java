package csw.chulbongkr.service.search;

import org.apache.lucene.analysis.TokenFilter;
import org.apache.lucene.analysis.TokenStream;
import org.apache.lucene.analysis.tokenattributes.CharTermAttribute;

import java.io.IOException;

import static csw.chulbongkr.util.LuceneUtil.extractInitialConsonants;

public class InitialConsonantTokenFilter extends TokenFilter {
    private final CharTermAttribute charTermAttribute = addAttribute(CharTermAttribute.class);

    protected InitialConsonantTokenFilter(TokenStream input) {
        super(input);
    }

    @Override
    public boolean incrementToken() throws IOException {
        if (input.incrementToken()) {
            String token = charTermAttribute.toString();
            String initialConsonants = extractInitialConsonants(token);
            charTermAttribute.setEmpty();
            charTermAttribute.append(initialConsonants);
            return true;
        }
        return false;
    }
}