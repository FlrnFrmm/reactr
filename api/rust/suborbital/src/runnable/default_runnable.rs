use super::{Runnable, Error};

/// A dummy struct to hold down the fort until a real Runnable is set
pub struct DefaultRunnable;

impl DefaultRunnable {
	pub const fn new() -> Self {
		Self
	}
}

impl Runnable for DefaultRunnable {
    fn run(&self, _input: Vec<u8>) -> Result<Vec<u8>, Error> {
        Err(Error::new_run_error(500, ""))
    }
}
