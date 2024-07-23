package csw.chulbongkr.controller.marker.search;

import csw.chulbongkr.entity.lucene.MarkerSearch;
import csw.chulbongkr.service.search.LuceneService;
import lombok.RequiredArgsConstructor;
import org.apache.lucene.queryparser.classic.ParseException;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

import java.io.IOException;
import java.util.List;

@RequiredArgsConstructor
@RestController
@RequestMapping("/api/v1/markers")
public class SearchController {
    private final LuceneService luceneService;

    @GetMapping("/search")
    public List<MarkerSearch> search(@RequestParam String term) throws IOException, ParseException {
        return luceneService.searchMarkers(term);
    }
}
