package config

import (
	op "imagecore/operation" // Grab `EncoderOption` from operation package.
)

// This ts a utility function to convert pipeline block to image operation.

func PipelineBlockToOperation(pb PipelineBlock) op.Operation {

	// We assume the loaded config has been fully checked.
	// Therefore we don't return error here.
	switch pb.Operation {

	case OperationDecode:
		return op.Decode()

	case OperationResize: // Resize block.
		if pb.Resize.Factor != 0.0 { // `Factor` nas first priority.
			return op.ResizeImageByFactor(pb.Resize.Algorithm, pb.Resize.Factor)
		} else if pb.Resize.Width != 0 { // If `Factor` is not set, use `Width`.
			return op.ResizeImageByWidth(pb.Resize.Algorithm, pb.Resize.Width)
		} else if pb.Resize.Height != 0 { // If `Width` is not set, use `Height`.
			return op.ResizeImageByHeight(pb.Resize.Algorithm, pb.Resize.Height)
		} else { // Empty resize block.
			return nil // This should not happen, since the config has been checked.
		}

	case OperationCrop: // Crop block.
		return op.Crop(pb.Crop.Width, pb.Crop.Height, pb.Crop.Alignment)

	case OperationEncode: // Encode block.
		return op.Encode(pb.Encode.Format, (*op.EncoderOption)(pb.Encode.Options))

	case OperationIccEmbed: // ICC embedding block.
		return op.EmbedProfile(pb.ICCEmbedProfile.ProfileName)

	case OperationWrite: // File output block.
		return op.WriteImageToFile(pb.Write.GenerateFileName())
	default:
		return nil // This should not happen, since the config has been checked.
	}
}
