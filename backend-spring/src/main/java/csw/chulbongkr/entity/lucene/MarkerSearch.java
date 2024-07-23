package csw.chulbongkr.entity.lucene;

import lombok.Data;

@Data
public class MarkerSearch {
    private int markerId;
    private String address;
    private String province;
    private String city;
    private String fullAddress;
    private String initialConsonants;
}
