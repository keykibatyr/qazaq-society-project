import express from "express";
import nodemailer from "nodemailer";
import bodyParser from "body-parser";
import dotenv from "dotenv";

dotenv.config(); // loads .env file for security

const app = express();
app.use(express.static("public"));
app.use(bodyParser.urlencoded({ extended: true }));
app.use(bodyParser.json());

// Contact form endpoint
app.post("/send", async (req, res) => {
  const { name, email, message } = req.body;

  try {
    const transporter = nodemailer.createTransport({
      host: "sandbox.smtp.mailtrap.io",
      port: 2525,
      auth: {
        user: process.env.MAILTRAP_USER,
        pass: process.env.MAILTRAP_PASS,
      },
    });

    await transporter.sendMail({
      from: `"Qazaq Society Website" <no-reply@qazaqsociety.be>`,
      to: "b.kapizov@gmail.com", // your test email
      subject: "📩 New Message from Qazaq Society Contact Form",
      text: `
Name: ${name}
Email: ${email}
Message:
${message}
      `,
    });

    console.log("✅ Email sent successfully!");
    res.status(200).send("Message sent successfully!");
  } catch (error) {
    console.error("❌ Error sending email:", error);
    res.status(500).send("Failed to send message. Please try again later.");
  }
});

app.listen(3000, () =>
  console.log("🚀 Server running at http://localhost:3000")
);
