/* client.java 
 * Name: Brendan Lin
 */

//Imports
import java.net.*;
import java.io.*;

public class Client {
  public static void main(String[] args) throws IOException {
	
    String host;
	int port;

	//Check command line args length
	if (args.length == 0) {
      host = "localhost";
	  port = 6666;
	} else {
	  host = args[0];
	  port = Integer.parseInt(args[1]);
	}

	//Connect to TCP and setup user IO
	try (
	  var s = new Socket(host, port); 
	  var in =
	    new BufferedReader(new InputStreamReader(s.getInputStream()));
	  var out =
	    new PrintWriter(new OutputStreamWriter(s.getOutputStream()), true);
	  var stdin =
	    new BufferedReader(new InputStreamReader(System.in));
	  ) {

	    //Start thread to listen to server for when messages are sent concurrently
	    new Thread(new ResponseListener(s)).start();

		//Run client
		String line;

		System.out.println("Enter commands below: ");

		while (true) {
		  if ((line = stdin.readLine()) == null) {
			break;
		  }

		  out.println(line);
		}
	 }
  }
}

//Thread class to get server responses concurrently
class ResponseListener implements Runnable {
  private Socket client;

  public ResponseListener(Socket client) {
	this.client = client;
  }

  @Override
  public void run() {
	try {
	  BufferedReader in = 
		  new BufferedReader(new InputStreamReader(client.getInputStream()));
	  String message;

	  while ((message = in.readLine()) != null) {
		System.out.println(message);
	  }

	  in.close();
	} catch (Exception e) {
	  System.out.println("Socket Closed");
	}
  }
}
