package utils

const markerPagination = `
SELECT 
    M.MarkerID, 
    M.UserID, 
    ST_Y(M.Location) AS Longitude, 
    ST_X(M.Location) AS Latitude, 
    M.Description, 
    U.Username, 
    M.CreatedAt, 
    M.UpdatedAt, 
    IFNULL(D.DislikeCount, 0) AS DislikeCount
FROM 
    Markers M
INNER JOIN 
    Users U ON M.UserID = U.UserID
LEFT JOIN 
    (
        SELECT 
            MarkerID, 
            COUNT(DislikeID) AS DislikeCount
        FROM 
            MarkerDislikes
        GROUP BY 
            MarkerID
    ) D ON M.MarkerID = D.MarkerID
WHERE 
    M.UserID = ?
ORDER BY 
    M.CreatedAt DESC
LIMIT ? OFFSET ?
`
