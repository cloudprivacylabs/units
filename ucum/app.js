const express = require('express');
const app = express();
const port = process.env.port || 8080;

app.use(express.json());

app.listen(port, err => {
    if (err) {
        return console.log("ERROR", err)
    }
    console.log(`Listening on port ${port}`);
})

// UCUM 
const ucum = require('@lhncbc/ucum-lhc');
const utils = ucum.UcumLhcUtils.getInstance();

app.get("/validate", (req, res) => {
    // handle UCUM validate function

    let unit = req.query.unit;
    let ret = utils.validateUnitString(unit);
    // if (ret['status'] === 'valid') {
    //     res.json(ret);
    //     res.end();
    // /* the conversion was successful.
    //     returnObj['toVal'] will contain the conversion result
    //     (~1943.9999999999998 - number, not formatted string)
    //     returnObj['msg'] will be null
    //     returnObj['fromUnit'] will contain the unit object for [fth_us]
    //     returnObj['toUnit'] will contain the unit object for [in_us]
    // */
    // } else {
    //     res.send(ret.msg);
    //     res.end();
    // /* the conversion encountered an error
    //     returnObj['toVal'] will be null
    //     returnObj['msg'] will contain a message describing the error
    //     returnObj['fromUnit'] will be null
    //     returnObj['toUnit'] will be null
    // */
    // }
    res.json(ret);
    console.log(ret);
})

app.get("/convert", (req, res) => {
    // handle UCUM convert function

    let fromUnitCode = req.query.unit;
    let toUnitCode = req.query.output;
    let value = req.query.value;
    let ret = utils.convertUnitTo(fromUnitCode, value, toUnitCode);
    // if (ret['status'] === 'succeeded') {
    //     res.json(ret);
    // /* the conversion was successful.
    //     returnObj['toVal'] will contain the conversion result
    //     (~1943.9999999999998 - number, not formatted string)
    //     returnObj['msg'] will be null
    //     returnObj['fromUnit'] will contain the unit object for [fth_us]
    //     returnObj['toUnit'] will contain the unit object for [in_us]
    // */
    // } else if (ret['status'] === 'failed') {
    //     res.send(ret.msg);
    // /* the conversion could not be made.
    //     returnObj['toVal'] will be null
    //     returnObj['msg'] will contain a message describing the failure
    //     returnObj['fromUnit'] will be null
    //     returnObj['toUnit'] will be null
    // */
    // } else if (ret['status'] === 'error') {
    //     res.status(404).end();
    // /* the conversion encountered an error
    //     returnObj['toVal'] will be null
    //     returnObj['msg'] will contain a message describing the error
    //     returnObj['fromUnit'] will be null
    //     returnObj['toUnit'] will be null
    // */
    // }
    res.json(ret);
    console.log(ret);
});
