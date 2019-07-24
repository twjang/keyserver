package com.terra.wallet;

import java.io.BufferedReader;
import java.io.DataOutputStream;
import java.io.IOException;
import java.io.InputStreamReader;
import java.net.HttpURLConnection;
import java.net.URL;

import org.json.JSONObject;

// import javax.net.ssl.HttpsURLConnection;


public class App 
{
    public static void main( String[] args )
    {
        final String URL = "http://127.0.0.1:3000";
        
        String sendTx = getSendTx(URL);
        if (sendTx == "") {
            return;
        }

        String signedTx = getSignedTx(URL, sendTx);
        System.out.println( signedTx);
        String response = broadcastTx(URL, signedTx);
        System.out.println( response);

    }

    static private String getSendTx(String urlStr) {
        try {
            URL url = new URL(urlStr + "/tx/bank/send");
            HttpURLConnection con = (HttpURLConnection) url.openConnection();
            
            con.setRequestMethod("POST");
            con.setDoOutput(true);


            JSONObject json = new JSONObject();
            json.put("sender", "terra1t849fxw7e8ney35mxemh4h3ayea4zf77dslwna");
            json.put("reciever", "terra1v9ku44wycfnsucez6fp085f5fsksp47u9x8jr4");
            json.put("amount", "1000000uluna");
            json.put("memo", "937767194");
            json.put("chain_id", "columbus-2");
            json.put("gas_adjustment", "1.4");
            json.put("gas_prices", "0.015ukrw");

            DataOutputStream wr = new DataOutputStream(con.getOutputStream()); 
            wr.writeBytes(json.toString()); 
            wr.flush(); 
            wr.close(); 

            int responseCode = con.getResponseCode(); 
            BufferedReader in = new BufferedReader(new InputStreamReader(con.getInputStream())); 
            String inputLine; StringBuffer response = new StringBuffer(); 
            
            while ((inputLine = in.readLine()) != null) { 
                response.append(inputLine); 
            } 
            
            in.close(); 

            if (responseCode == 200) {
                return response.toString();
            }

            System.out.println("Failed to Get Msg; Status Code " + responseCode);
            System.out.println("Response: " + response);
            
        } catch (IOException ex) {
            System.out.println(ex);
        }

        return "";
    }

    static private String getSignedTx(String urlStr, String tx) {
        try {
            URL url = new URL(urlStr + "/tx/sign");
            HttpURLConnection con = (HttpURLConnection) url.openConnection();
            
            con.setRequestMethod("POST");
            con.setDoOutput(true);


            JSONObject json = new JSONObject();
            json.put("tx",  new JSONObject(tx));
            json.put("name", "tmp");
            json.put("passphrase", "12345678");
            json.put("chain_id", "columbus-2");
            json.put("account_number", "93");
            json.put("sequence", "64");

            DataOutputStream wr = new DataOutputStream(con.getOutputStream()); 
            wr.writeBytes(json.toString()); 
            wr.flush(); 
            wr.close(); 

            int responseCode = con.getResponseCode(); 
            BufferedReader in = new BufferedReader(new InputStreamReader(con.getInputStream())); 
            String inputLine; StringBuffer response = new StringBuffer(); 
            
            while ((inputLine = in.readLine()) != null) { 
                response.append(inputLine); 
            } 
            
            in.close(); 

            if (responseCode == 200) {
                return response.toString();
            }

            System.out.println("Failed to Get Msg; Status Code " + responseCode);
            System.out.println("Response: " + response);
            
        } catch (IOException ex) {
            System.out.println(ex.getMessage());
        }

        return "";
    }

    static private String broadcastTx(String urlStr, String signedTx) {
        try {
            URL url = new URL(urlStr + "/tx/broadcast");
            HttpURLConnection con = (HttpURLConnection) url.openConnection();
            
            con.setRequestMethod("POST");
            con.setDoOutput(true);


            JSONObject json = new JSONObject(signedTx);

            DataOutputStream wr = new DataOutputStream(con.getOutputStream()); 
            wr.writeBytes(json.toString()); 
            wr.flush(); 
            wr.close(); 

            int responseCode = con.getResponseCode(); 
            BufferedReader in = new BufferedReader(new InputStreamReader(con.getInputStream())); 
            String inputLine; StringBuffer response = new StringBuffer(); 
            
            while ((inputLine = in.readLine()) != null) { 
                response.append(inputLine); 
            } 
            
            in.close(); 

            if (responseCode == 200) {
                return response.toString();
            }

            System.out.println("Failed to Get Msg; Status Code " + responseCode);
            System.out.println("Response: " + response);
            
        } catch (IOException ex) {
            System.out.println(ex.getMessage());
        }

        return "";
    }
}