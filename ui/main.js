function fetchData() {
    console.log("fetchData called"); // add logging

    // Clear the existing data
    let resultsTable = document.getElementById('results-table');
    console.log("Existing rows: ", resultsTable.rows.length); // add logging
    while (resultsTable.rows.length > 1) {
        resultsTable.deleteRow(1);
    }

    // Get the origin airport
    let origin = document.getElementById('origin').value;
    console.log("Origin: ", origin); // add logging

    // Send GET request to our server
    fetch('/api/get-destinations/?origin=' + origin)
        .then(response => {
            console.log("Server Response: ", response); // add logging
            return response.json();
        })
        .then(data => {
            console.log("Data received: ", data); // add logging
            data.data.forEach(flight => {
                // Create new row
                let row = resultsTable.insertRow();

                // Create cells in row
                let destinationCell = row.insertCell();
                let departureDateCell = row.insertCell();
                let priceCell = row.insertCell();
                let bookLinkCell = row.insertCell();

                destinationCell.innerHTML = flight.destination;
                departureDateCell.innerHTML = flight.departureDate;
                priceCell.innerHTML = flight.price.total;
                bookLinkCell.innerHTML = '<a href="' + flight.links.flightOffers + '" target="_blank">Book Now</a>';

                console.log("Added new row for destination: ", flight.destination); // add logging
            });
        })
        .catch((error) => {
            console.error('Error:', error);  // catch any error and log it
        });
}

