package csw.chulbongkr.service.search;

import csw.chulbongkr.entity.lucene.MarkerSearch;
import jakarta.annotation.PostConstruct;
import lombok.extern.slf4j.Slf4j;
import org.apache.lucene.analysis.Analyzer;
import org.apache.lucene.document.Document;
import org.apache.lucene.document.Field;
import org.apache.lucene.document.TextField;
import org.apache.lucene.index.DirectoryReader;
import org.apache.lucene.index.IndexWriter;
import org.apache.lucene.index.IndexWriterConfig;
import org.apache.lucene.index.Term;
import org.apache.lucene.queryparser.classic.ParseException;
import org.apache.lucene.search.*;
import org.apache.lucene.store.Directory;
import org.apache.lucene.store.MMapDirectory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.io.IOException;
import java.nio.file.DirectoryStream;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.List;

import static csw.chulbongkr.util.LuceneUtil.*;

@Slf4j
@Service
public class LuceneService {

    private Directory indexDirectory;
    private Analyzer analyzer;
    private IndexWriter indexWriter;

    @Autowired
    private CustomAnalyzerProvider customAnalyzerProvider;

    @PostConstruct
    public void init() throws IOException {
        Path indexPath = Paths.get("target/lucene/index");
        cleanDirectory(indexPath);  // Clean the index directory before initializing

        indexDirectory = new MMapDirectory(indexPath);
        analyzer = customAnalyzerProvider.getAnalyzer();
        IndexWriterConfig config = new IndexWriterConfig(analyzer);
        indexWriter = new IndexWriter(indexDirectory, config);
    }

    public void indexMarkerBatch(List<MarkerSearch> markers) throws IOException {
        List<Document> documents = new ArrayList<>();
        for (MarkerSearch marker : markers) {
            Document doc = new Document();
            doc.add(new TextField("markerId", String.valueOf(marker.getMarkerId()), Field.Store.YES));
            doc.add(new TextField("address", marker.getAddress(), Field.Store.YES));
            doc.add(new TextField("province", marker.getProvince(), Field.Store.YES));
            doc.add(new TextField("city", marker.getCity(), Field.Store.YES));
            doc.add(new TextField("fullAddress", marker.getFullAddress(), Field.Store.YES));
            doc.add(new TextField("initialConsonants", marker.getInitialConsonants(), Field.Store.YES));
            documents.add(doc);
        }
        for (Document doc : documents) {
            indexWriter.addDocument(doc);
        }
        indexWriter.commit();
    }

    private void cleanDirectory(Path path) throws IOException {
        if (Files.exists(path) && Files.isDirectory(path)) {
            try (DirectoryStream<Path> directoryStream = Files.newDirectoryStream(path)) {
                for (Path file : directoryStream) {
                    Files.delete(file);
                }
            }
        } else {
            Files.createDirectories(path);
        }
    }

    public void indexMarker(MarkerSearch marker) throws IOException {
        Document doc = new Document();
        doc.add(new TextField("markerId", String.valueOf(marker.getMarkerId()), Field.Store.YES));
        doc.add(new TextField("address", marker.getAddress(), Field.Store.YES));
        doc.add(new TextField("province", marker.getProvince(), Field.Store.YES));
        doc.add(new TextField("city", marker.getCity(), Field.Store.YES));
        doc.add(new TextField("fullAddress", marker.getFullAddress(), Field.Store.YES));
        doc.add(new TextField("initialConsonants", marker.getInitialConsonants(), Field.Store.YES));
        indexWriter.addDocument(doc);
        indexWriter.commit();
    }

    public List<MarkerSearch> searchMarkers(String queryStr) throws IOException, ParseException {
        List<MarkerSearch> results = new ArrayList<>();
        DirectoryReader reader = DirectoryReader.open(indexDirectory);
        IndexSearcher searcher = new IndexSearcher(reader);

        List<Query> queries = new ArrayList<>();

        // Higher boost for full queryStr matches
        queries.add(new BoostQuery(new TermQuery(new Term("fullAddress", queryStr)), 25.0f));
        queries.add(new BoostQuery(new PhraseQuery("fullAddress", queryStr), 20.0f));
        queries.add(new BoostQuery(new WildcardQuery(new Term("fullAddress", "*" + queryStr + "*")), 15.0f));
        queries.add(new BoostQuery(new PrefixQuery(new Term("fullAddress", queryStr)), 35.0f));

        String brokenConsonants = extractInitialConsonants(segmentConsonants(queryStr));
        queries.add(new BoostQuery(new TermQuery(new Term("initialConsonants", brokenConsonants)), 15.0f));
        queries.add(new BoostQuery(new WildcardQuery(new Term("initialConsonants", "*" + brokenConsonants + "*")), 7.0f));
        queries.add(new BoostQuery(new PrefixQuery(new Term("initialConsonants", brokenConsonants)), 25.0f));

        String standardizedProvince = standardizeProvince(queryStr);
        if (!standardizedProvince.equals(queryStr)) {
            queries.add(new BoostQuery(new PrefixQuery(new Term("province", standardizedProvince)), 3.0f));
        } else {
            queries.add(new BoostQuery(new PrefixQuery(new Term("city", queryStr)), 10.0f));
            queries.add(new BoostQuery(new TermQuery(new Term("city", queryStr)), 10.0f));
            queries.add(new BoostQuery(new TermQuery(new Term("district", queryStr)), 5.0f));
            queries.add(new BoostQuery(new PrefixQuery(new Term("address", queryStr)), 10.0f));
            queries.add(new BoostQuery(new WildcardQuery(new Term("address", "*" + queryStr + "*")), 5.0f));
        }

        // Split queryStr by spaces and create additional queries
        String[] terms = queryStr.split("\\s+");
        for (String term : terms) {
            if (!term.isEmpty()) {
                queries.add(new BoostQuery(new TermQuery(new Term("fullAddress", term)), 5.0f));
                queries.add(new BoostQuery(new WildcardQuery(new Term("fullAddress", "*" + term + "*")), 3.0f));
                queries.add(new BoostQuery(new PrefixQuery(new Term("fullAddress", term)), 7.0f));

                String brokenTermConsonants = extractInitialConsonants(segmentConsonants(term));
                queries.add(new BoostQuery(new TermQuery(new Term("initialConsonants", brokenTermConsonants)), 5.0f));
                queries.add(new BoostQuery(new WildcardQuery(new Term("initialConsonants", brokenTermConsonants + "*")), 3.0f));
                queries.add(new BoostQuery(new WildcardQuery(new Term("initialConsonants", "*" + brokenTermConsonants + "*")), 1.0f));
                queries.add(new BoostQuery(new PrefixQuery(new Term("initialConsonants", brokenTermConsonants)), 7.0f));

                queries.add(new BoostQuery(new PrefixQuery(new Term("city", term)), 3.0f));
                queries.add(new BoostQuery(new TermQuery(new Term("city", term)), 3.0f));
                queries.add(new BoostQuery(new TermQuery(new Term("district", term)), 2.0f));
                queries.add(new BoostQuery(new PrefixQuery(new Term("address", term)), 3.0f));
                queries.add(new BoostQuery(new WildcardQuery(new Term("address", "*" + term + "*")), 1.0f));
            }
        }

        BooleanQuery.Builder booleanQuery = new BooleanQuery.Builder();
        for (Query query : queries) {
            booleanQuery.add(query, BooleanClause.Occur.SHOULD);
        }

        TopDocs docs = searcher.search(booleanQuery.build(), 10);
        for (ScoreDoc scoreDoc : docs.scoreDocs) {
            Document doc = searcher.storedFields().document(scoreDoc.doc);
            MarkerSearch marker = new MarkerSearch();
            marker.setMarkerId(Integer.parseInt(doc.get("markerId")));
            marker.setAddress(doc.get("address"));
            marker.setProvince(doc.get("province"));
            marker.setCity(doc.get("city"));
            marker.setFullAddress(doc.get("fullAddress"));
            marker.setInitialConsonants(doc.get("initialConsonants"));
            results.add(marker);
        }

        reader.close();
        return results;
    }
}
