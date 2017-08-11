# go_lan_file_server
A simple implementation of a LAN file server that supports file uploads that can resume from where it left off.

---

### The Program
This was a small project to create a local file server that I could write an accompanying phone app for to push files(photos and videos mostly) to my computer while I am at my house. The other thing I wanted to allow this program to do was have a feature to allow you to pick back up uploading a file where you left off. I created a simple file format named "SAVE" to keep track of uploaded files and handle the upload resuming feature. The file server in its current state is very simple and generic. I have made it to where you can save header data with the files you upload and you can define the header on the client side. The downside right now is you cannot change from the header format you initially start off with. However, the Header is an interface and you can implement different header logic if you need something else that is not so simple. 

### The Files
- Main.go - the main file; the paths and their logic are defined here.
- src/sfile/sfile.go - the file that implements the SAVE file format logic and the associated objects and interfaces.
- src/sfile/sheader.go - imlpements a SimpleHeader object that adheres to the HeaderFormat interface. This object is for very simple uses.

### Current Paths
#### /post_file - POST request 
- takes json format:
  - Data - base64 encoded byte array of file data.
  - ValidateFile - base64 encoded byte array of sha256 value of file data.
  - StartIndex - integer of starting position of range of file data you are sending.
  - Size - integer of size of your entire file.
  - Attributes - map[string]string. This is your custom header format.
- returns json format:
  - Error - empty if everything is okay, message if not.
  - Count - integer, 0 if Error is set, but if Error is set and this is greater than 0 then the file already exists and the value for Count is the last position in your file you stopped at.
#### /get_folders GET request 
- takes nothing.
- returns json format:
  - Folders - array of folder objects that have keys "Name"(folder name) and "Count"(How many files in folder).
  - Error - empty if nothing wrong, message otherwise.
#### /get_files POST request
- takes json format:
  - Folder - string, The folder you want to pull files from.
  - StartIndex - integer, of which file you want to start grabbing from. 0 based index.
  - EndIndex - integer, of the last position(exclusively) of the files you would like to grab.
  - Attributes - map[string]string, You only need to set the keys of the attributes for your header format so it can pull and set them to the right keys when returned.
- returns json format:
  - Same as what /post_file takes as a json format
#### /validate_file GET request
- takes GET parameters.
  - Folder - string, The folder to reference the file from.
  - Index or Hash
    - Index - int, The index of the file from the list of files in the folder. This uses the initial stored hash of the indexed file to compare against the sha256 hash of the stored data from the indexed file.
    - Hash - string, The sha256 hash of the file. This will see if a file with this hash already exists in the folder.

### Features still needed
- update README to document how to validate file with sha256
- edit file paths to use generic separator so this can work on more than just windows
