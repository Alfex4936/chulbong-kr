package csw.chulbongkr.config;

import csw.chulbongkr.service.search.LuceneIndexerService;
import org.springframework.stereotype.Component;

@Component
public class Initializer {

    private final LuceneIndexerService luceneIndexerService;

    public Initializer(LuceneIndexerService luceneIndexerService) {
        this.luceneIndexerService = luceneIndexerService;
        // Explicitly call a method on LuceneIndexerService to ensure it's initialized
        luceneIndexerService.ensureInitialized();
    }
}