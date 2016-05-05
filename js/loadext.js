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
        // send the object back to and load in HTML
        loadList(json_list);
    };
    fr.readAsText(file);

    // change visual html label to filename
    $('#label_file_input').text(file['name']);

}

function parseList(str) {
    //var regexp = /^(.*):(.*)$/;
    // get only rows with a correct format
    var regexp = /(.*)\s*:\s*(.*)/g;
    var rows = str.match(regexp); //array

    var dishes = [];
    // loop, extract elements
    rows.forEach(function(row) {
        regexp.lastIndex = 0; // to be able to reuse the same regexp in the loop
        var rowdata = regexp.exec(row); //1:name, 2:comma-separated ingredients
        dishes.push({"name":rowdata[1],"ingredients":rowdata[2].split(/\s*,\s*/)});
    });
    return {"dishes":dishes};
}
