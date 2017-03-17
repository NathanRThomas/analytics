package anal

  //-------------------------------------------------------------------------//
 //----- IMPORTS -----------------------------------------------------------//
//-------------------------------------------------------------------------//

import (
    "bytes"
    "fmt"
    "net/http"
//    "strconv"
    
    "github.com/ajg/form"
    )

  //-------------------------------------------------------------------------//
 //----- CONSTANTS ---------------------------------------------------------//
//-------------------------------------------------------------------------//

const (
    anal_tracking_type_ga   = iota
    )

const anal_ga_endpoint = "https://www.google-analytics.com/collect"

  //-------------------------------------------------------------------------//
 //----- STRUCTS -----------------------------------------------------------//
//-------------------------------------------------------------------------//

type gaParams struct {
    Version     int     `form:"v"`
    TrackingID  string  `form:"tid"`
    UserID      string  `form:"uid"`
    HitType     string  `form:"t"`
    DataSource  string  `form:"ds,omitempty"`
    Location    string  `form:"geoid,omitempty"`
    Category    string  `form:"ec,omitempty"`
    Action      string  `form:"ea,omitempty"`
    Label       string  `form:"el,omitempty"`
    Value       int     `form:"ev,omitempty"`
    Page        string  `form:"dl,omitempty"`
    Screen      string  `form:"cd,omitempty"`
    App         string  `form:"an,omitempty"`
    Title       string  `form:"dt,omitempty"`
}

type Client struct {
    trackingType    int
    trackingID      string
}


  //-------------------------------------------------------------------------//
 //----- PRIVATE FUNCTIONS -------------------------------------------------//
//-------------------------------------------------------------------------//

/*! \brief Final step once we have the params body created
 */
func (c Client) postGA (params gaParams) (error) {
    //set some defaults
    params.Version = 1
    params.TrackingID = c.trackingID
    //params.App = "Headliner Labs"
    
    str, err := form.EncodeToString(params)
    //fmt.Println(str)
    if err != nil { return err }    //shouldn't happen
    
    buf := bytes.NewBufferString(str)

    req, err := http.NewRequest("POST", anal_ga_endpoint, buf)

    if err == nil {
        req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
        //req.Header.Set("Content-Length", strconv.Itoa(len(buf)))
        
        client := &http.Client{}
        resp, err := client.Do(req)
        //fmt.Printf("%+v\n", resp)
        if err != nil {
            return fmt.Errorf("Google Analytics client Failed: " + err.Error())
        }
        defer resp.Body.Close() //we're done with this request
        return nil
    } else {
        return fmt.Errorf("Google Analytics client Failed: " + err.Error())
    }
}


  //-------------------------------------------------------------------------//
 //----- PUBLIC FUNCTIONS --------------------------------------------------//
//-------------------------------------------------------------------------//

func (c Client) Event (clientID, countryCode, referrer, category, action, label string, value int) (error) {
    //now figure out which analytics backend we're using
    switch c.trackingType {
        case anal_tracking_type_ga:
            params := gaParams {UserID: clientID, HitType: "event", DataSource: referrer, Location: countryCode,
                Category: category, Action: action, Label: label, Value: value}
            
            return c.postGA (params)
    }
    
    return nil
}

func (c Client) Page (clientID, countryCode, referrer, page string) (error) {
    //now figure out which analytics backend we're using
    switch c.trackingType {
        case anal_tracking_type_ga:
            params := gaParams {UserID: clientID, HitType: "pageview", DataSource: referrer, Location: countryCode,
                Page: page, Title: page}
            return c.postGA (params)
    }
    return nil
}

func (c Client) Screen (clientID, countryCode, referrer, screen string) (error) {
    //now figure out which analytics backend we're using
    switch c.trackingType {
        case anal_tracking_type_ga:
            params := gaParams {UserID: clientID, HitType: "screenview", DataSource: referrer, Location: countryCode,
                Screen: screen, Title: screen}
            //fmt.Printf("%+v\n", params)
            return c.postGA (params)
    }
    return nil
}