package utils

const getAllMarkers = `
SELECT 
M.MarkerID, 
M.UserID, 
ST_Y(M.Location) AS Longitude, 
ST_X(M.Location) AS Latitude, 
M.Description, 
COALESCE(U.Username, '탈퇴한 사용자') AS Username, 
M.CreatedAt, 
M.UpdatedAt, 
IFNULL(D.DislikeCount, 0) AS DislikeCount
FROM 
Markers M
LEFT JOIN 
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
) D ON M.MarkerID = D.MarkerID;`

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
