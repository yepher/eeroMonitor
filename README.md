## Eero Monitor


## Usage

* If you do not already have a session key you will need to generate one. This is a two step process:
	* Request verification code:
		* `eeroMonitor --loginID=[email@example.com]`
		* This returns a session key and causes Eero to send a verification code via email
	* Verify code to activate session key:
		* `eeroMonitor -verificationKey=123456 -sessionKey="6654321|232c2aoj93fvdes82eg99ase7e"`


* Once the session key is is verified pass that as `-sessionKey`