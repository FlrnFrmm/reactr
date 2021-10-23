#[derive(Debug)]
pub enum Error {
	Run { code: i32, message: String },
	Host { message: String }
}

impl Error {
	pub fn new_host_error<T: Into<String>>(message: T) -> Self {
		Self::Host { message: message.into() }
	}
	pub fn new_run_error<T: Into<String>>(code: i32, message: T) -> Self {
		Self::Run { code, message: message.into() }
	}
}

impl std::fmt::Display for Error {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
			Self::Host { message} => write!(f, "Host Error: {}", message),
			Self::Run { code, message } => write!(f, "Run Error({}): {}", code, message)
		}
    }
}

impl std::error::Error for Error {}