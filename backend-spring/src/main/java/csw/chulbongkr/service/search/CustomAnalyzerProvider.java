package csw.chulbongkr.service.search;

import org.apache.lucene.analysis.Analyzer;
import org.apache.lucene.queryparser.classic.ParseException;
import org.springframework.stereotype.Component;

import java.io.IOException;

@Component
public class CustomAnalyzerProvider {

    private final Analyzer customAnalyzer;

    public CustomAnalyzerProvider() throws IOException, ParseException {
        this.customAnalyzer = new CustomKoreanAnalyzer();
    }

    public Analyzer getAnalyzer() {
        return customAnalyzer;
    }
}