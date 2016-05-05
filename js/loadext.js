/* loadext.js */

function loadExternalList(obj) {
    // abort if not supported
    if (!window.File || !window.FileReader || !window.FileList || !window.Blob) {
        console.log("File API not supported");
        return;
    }

    // get file
    var file;
    file = obj.files[0];

    // get data using the obviously supported FileReader
    var fr;
    fr = new FileReader();
    fr.onload = function() {
        // parse the data into a JSON object
        var json_list = parseList(fr.result);
        if (json_list !== false) {
            // load in HTML
            loadList(json_list);
            // change visual html label to filename
            $('#label_file_input').text(file['name']);
        } else {
            alert("ERROR: Could not find any valid rows in file.");
        }
    };
    fr.readAsText(file);
}

function parseList(str) {
    //var regexp = /^(.*):(.*)$/;
    // get only rows with a correct format
    var regexp = /(.*)\s*:\s*(.*)/g;
    var rows;
    rows = str.match(regexp); //array

    // if no rows had correct format, abort
    if (rows === null) return false;

    var dishes = [];
    // loop, extract elements
    rows.forEach(function(row) {
        regexp.lastIndex = 0; // to be able to reuse the same regexp in the loop
        var rowdata = regexp.exec(row); //1:name, 2:comma-separated ingredients
        dishes.push({"name":rowdata[1],"ingredients":rowdata[2].split(/\s*,\s*/)});
    });
    return {"dishes":dishes};
}
